/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: MeteorShower.tsx
Description: Animated meteor shower background for SteriaFront. Uses React and canvas for a magical, modern effect.
*/

import React, { useRef, useEffect } from 'react';

interface Meteor {
  x: number;
  y: number;
  length: number;
  speed: number;
  angle: number;
  opacity: number;
  width: number;
  color: string;
}

const METEOR_COLORS = [
  'rgba(244,114,182,0.7)', // pink
  'rgba(162,28,175,0.7)',  // purple
  'rgba(99,102,241,0.7)',  // indigo
  'rgba(236,72,153,0.7)',  // fuchsia
  'rgba(255,255,255,0.5)', // white
];

const METEOR_COUNT = 18;

const randomMeteor = (w: number, h: number): Meteor => ({
  x: Math.random() * w,
  y: Math.random() * h * 0.7,
  length: 80 + Math.random() * 60,
  speed: 2.5 + Math.random() * 2.5,
  angle: Math.PI / 4 + (Math.random() - 0.5) * 0.2,
  opacity: 0.5 + Math.random() * 0.5,
  width: 2 + Math.random() * 2,
  color: METEOR_COLORS[Math.floor(Math.random() * METEOR_COLORS.length)],
});

const MeteorShower: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const meteors = useRef<Meteor[]>([]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let animationId: number;
    let width = window.innerWidth;
    let height = window.innerHeight;

    // Resize canvas to fill window
    const resize = () => {
      width = window.innerWidth;
      height = window.innerHeight;
      canvas.width = width;
      canvas.height = height;
    };
    resize();
    window.addEventListener('resize', resize);

    // Initialize meteors
    meteors.current = Array.from({ length: METEOR_COUNT }, () => randomMeteor(width, height));

    // Animation loop
    const animate = () => {
      ctx.clearRect(0, 0, width, height);
      for (const meteor of meteors.current) {
        ctx.save();
        ctx.globalAlpha = meteor.opacity;
        ctx.strokeStyle = meteor.color;
        ctx.lineWidth = meteor.width;
        ctx.beginPath();
        ctx.moveTo(meteor.x, meteor.y);
        ctx.lineTo(
          meteor.x - Math.cos(meteor.angle) * meteor.length,
          meteor.y - Math.sin(meteor.angle) * meteor.length
        );
        ctx.shadowColor = meteor.color;
        ctx.shadowBlur = 16;
        ctx.stroke();
        ctx.restore();

        // Move meteor
        meteor.x += Math.cos(meteor.angle) * meteor.speed;
        meteor.y += Math.sin(meteor.angle) * meteor.speed;

        // Respawn if out of bounds
        if (
          meteor.x < -meteor.length ||
          meteor.y > height + meteor.length ||
          meteor.x > width + meteor.length
        ) {
          Object.assign(meteor, randomMeteor(width, height));
          meteor.x = Math.random() * width;
          meteor.y = -20;
        }
      }
      animationId = requestAnimationFrame(animate);
    };
    animate();

    return () => {
      window.removeEventListener('resize', resize);
      cancelAnimationFrame(animationId);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        width: '100vw',
        height: '100vh',
        zIndex: 0,
        pointerEvents: 'none',
      }}
      aria-hidden="true"
    />
  );
};

export default MeteorShower; 