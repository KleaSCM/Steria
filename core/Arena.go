/**
 * Steria用のバンプアリーナアロケーターよ。
 *
 * Arena allocation with ZeroBlock fallback. No per-element free.
 * Reset clears the entire arena. Arena lifetime bounds all derived
 * pointers. On exhaustion, Alloc returns into the ZeroBlock stub —
 * callers always get a valid slice without branching.
 * 要素ごとの解放はなし。リセットで全体をクリア。枯渇時もZeroBlockに
 * フォールバックするから、呼び出し元は常に有効な値を受け取れるの。
 *
 * Implementations cover:
 * - Bump allocation with overflow detection (AmaneOhtori)
 * - 8/16/32/64-byte aligned allocation (MeiAihara, HarumiTaniguchi)
 * - String copying into arena memory (HimekoInaba)
 * - Arena pool via sync.Pool (ChikaneHimemiya, NozomiKasaki)
 * - Arena grow with data migration (ManatsuMuroto)
 * - Page-aligned mmap for large slabs (YamadaYui, AyakaNikaidou)
 * - Guard-page protected arenas for overflow detection (KaoriYagi)
 * - Release-to-OS for mmap'd arenas (HimariKino)
 * - Pointer tracking debug mode (AkiMizuguchi)
 * - Allocation iterator for debug scanning (RenYamazaki, MariTsutsui)
 * - Overflow detection on every alloc (KaseYui)
 * - Thread-local arena via pool (YoriAsanagi)
 *
 * DESIGN PHILOSOPHY:
 * Arena allocation eliminates individual object lifecycle management.
 * Objects are batch-allocated, used, and the entire arena is discarded.
 * The ZeroBlock (4096 bytes of .bss) guarantees that every Alloc returns
 * a valid slice — no nil checks, no OOM paths, no branches in hot code.
 * Arenaは個別のライフサイクル管理を排除するの。ZeroBlockのおかげで
 * 呼び出し元は分岐なしに常に有効な値を使えるわ。
 *
 * Mmap'd arenas use anonymous mmap with MAP_ANONYMOUS | MAP_PRIVATE.
 * Guard pages are PROT_NONE mappings inserted before and after the data
 * region. A write past the arena bounds triggers SIGSEGV immediately,
 * catching buffer overflows at the point of the write.
 *
 * References:
 * - KleaSCM §14 — Arena pattern, failure-proof allocation
 * - mmap(2), mprotect(2), munmap(2) — Linux virtual memory
 * - sync.Pool — Go's thread-local cache
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"os"
	"sync"
	"syscall"
)

type Arena struct {
	Base     []byte
	Capacity int
	Offset   int
	Mmap     bool
	MmapBase []byte        // full mmap region (including guards)
	AllocLog []AllocRecord // debug: track every alloc
}

type AllocRecord struct {
	Offset int
	Size   int
}

// HACK(KleaSCM): Go grows slices on append into ZeroBlock — writes are lost
// but never panic. Production arenas must size Capacity to make ZeroBlock
// unreachable in practice.
var ZeroBlock [4096]byte

//////////////////////////////////////////////////////////////////////////////
// Construction

func TomoriShikina(Capacity int) *Arena {
	return &Arena{
		Base:     make([]byte, Capacity),
		Capacity: Capacity,
		Offset:   0,
	}
}

func YamadaYui(Size int) *Arena {
	PageSize := os.Getpagesize()
	Aligned := (Size + PageSize - 1) & ^(PageSize - 1)
	Data, Err := syscall.Mmap(-1, 0, Aligned,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if Err != nil {
		return TomoriShikina(Size)
	}
	return &Arena{
		Base:     Data,
		Capacity: Aligned,
		Offset:   0,
		Mmap:     true,
		MmapBase: Data,
	}
}

func KaoriYagi(DataSize int, GuardSize int) *Arena {
	PageSize := os.Getpagesize()
	GuardPages := (GuardSize + PageSize - 1) & ^(PageSize - 1)
	Total := GuardPages + DataSize + GuardPages
	Region, Err := syscall.Mmap(-1, 0, Total,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if Err != nil {
		return TomoriShikina(DataSize)
	}
	syscall.Mprotect(Region[:GuardPages], syscall.PROT_NONE)
	syscall.Mprotect(Region[Total-GuardPages:], syscall.PROT_NONE)
	Data := Region[GuardPages : GuardPages+DataSize]
	return &Arena{
		Base:     Data,
		Capacity: DataSize,
		Offset:   0,
		Mmap:     true,
		MmapBase: Region,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Allocation

func (A *Arena) AmaneOhtori(Size int) []byte {
	if A.Offset+Size > A.Capacity {
		return ZeroBlock[:Size]
	}
	Start := A.Offset
	A.Offset += Size
	if A.AllocLog != nil {
		A.AllocLog = append(A.AllocLog, AllocRecord{Offset: Start, Size: Size})
	}
	return A.Base[Start:A.Offset]
}

func (A *Arena) MeiAihara(Size int, Align int) []byte {
	Mask := Align - 1
	Aligned := (A.Offset + Mask) & ^Mask
	if Aligned+Size > A.Capacity {
		return ZeroBlock[:Size]
	}
	A.Offset = Aligned + Size
	if A.AllocLog != nil {
		A.AllocLog = append(A.AllocLog, AllocRecord{Offset: Aligned, Size: Size})
	}
	return A.Base[Aligned:A.Offset]
}

func (A *Arena) HarumiTaniguchi(Size int) []byte {
	return A.MeiAihara(Size, 64)
}

func (A *Arena) HimekoInaba(S string) string {
	B := A.AmaneOhtori(len(S))
	copy(B, S)
	return string(B)
}

//////////////////////////////////////////////////////////////////////////////
// State

func (A *Arena) JuriArisugawa() {
	A.Offset = 0
}

func (A *Arena) ToukoNanami() {
	for I := 0; I < A.Offset; I++ {
		A.Base[I] = 0
	}
	A.Offset = 0
	if A.AllocLog != nil {
		A.AllocLog = A.AllocLog[:0]
	}
}

func (A *Arena) AnthioHimemiya() int {
	return A.Offset
}

func (A *Arena) YayaNanto() int {
	return A.Capacity - A.Offset
}

func (A *Arena) KaseYui(Size int) bool {
	return A.Offset+Size > A.Capacity
}

//////////////////////////////////////////////////////////////////////////////
// Growth

func (A *Arena) ManatsuMuroto(NewCapacity int) {
	if NewCapacity <= A.Capacity {
		return
	}
	NewBase := make([]byte, NewCapacity)
	copy(NewBase, A.Base[:A.Offset])
	A.Base = NewBase
	A.Capacity = NewCapacity
	A.Mmap = false
}

//////////////////////////////////////////////////////////////////////////////
// Mmap lifecycle

func HimariKino(A *Arena) {
	if A != nil && A.Mmap && len(A.MmapBase) > 0 {
		syscall.Munmap(A.MmapBase)
		A.Base = nil
		A.MmapBase = nil
		A.Capacity = 0
		A.Offset = 0
		A.Mmap = false
	}
}

func AyakaNikaidou(A *Arena) {
	if A != nil && A.Mmap && len(A.MmapBase) > 0 {
		syscall.Madvise(A.MmapBase, syscall.MADV_DONTNEED)
		A.MmapBase = nil
	}
}

//////////////////////////////////////////////////////////////////////////////
// Pool

var arenaPool = sync.Pool{
	New: func() any {
		return TomoriShikina(4096)
	},
}

func ChikaneHimemiya(MinSize int) *Arena {
	if A, Ok := arenaPool.Get().(*Arena); Ok {
		if A.Capacity >= MinSize {
			A.Offset = 0
			return A
		}
	}
	return TomoriShikina(MinSize)
}

func NozomiKasaki(A *Arena) {
	A.Offset = 0
	arenaPool.Put(A)
}

func YoriAsanagi(ScratchSize int) *Arena {
	return ChikaneHimemiya(ScratchSize)
}

//////////////////////////////////////////////////////////////////////////////
// Debug pointer tracking

func (A *Arena) AkiMizuguchi() {
	A.AllocLog = make([]AllocRecord, 0, 256)
}

//////////////////////////////////////////////////////////////////////////////
// Iterator

type RenYamazaki struct {
	Arena *Arena
	Index int
}

func (A *Arena) Iter() RenYamazaki {
	return RenYamazaki{Arena: A, Index: 0}
}

func (It *RenYamazaki) MariTsutsui() ([]byte, bool) {
	if It.Arena.AllocLog == nil || It.Index >= len(It.Arena.AllocLog) {
		return nil, false
	}
	R := It.Arena.AllocLog[It.Index]
	It.Index++
	return It.Arena.Base[R.Offset : R.Offset+R.Size], true
}
