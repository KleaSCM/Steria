/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Navbar.tsx
Description: Modern, animated, fully responsive Navbar for SteriaFront. Logo links to profile. Hamburger menu shows Logout on profile page. Uses Framer Motion, React Router Link, and modular SCSS.
*/

import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import styles from './Navbar.module.scss';

const navLinks = [
  { label: 'Home', href: '#' },
  { label: 'Explore', href: '#' },
  { label: 'Projects', href: '#' },
];

const Navbar: React.FC = () => {
  const [menuOpen, setMenuOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const onLogout = () => {
    setMenuOpen(false);
    // Clear session state here if needed
    navigate('/login');
  };

  const isProfile = location.pathname === '/profile';

  return (
    <nav className={styles.navbar}>
      <div className={styles.brand}>
        <Link to="/profile" className={styles.logo} tabIndex={0}>
          SteriaFront
        </Link>
      </div>
      <ul className={styles.links}>
        {navLinks.map(link => (
          <li key={link.label}>
            <a href={link.href}>{link.label}</a>
          </li>
        ))}
      </ul>
      <div className={styles.actions}>
        {!isProfile && (
          <>
            <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
              <Link to="/signup" className={styles.signUpBtn} tabIndex={0}>
                Sign Up
              </Link>
            </motion.div>
            <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
              <Link to="/login" className={styles.signInBtn} tabIndex={0}>
                Sign In
              </Link>
            </motion.div>
          </>
        )}
        {isProfile && (
          <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
            <button className={styles.signInBtn} onClick={onLogout} tabIndex={0}>
              Logout
            </button>
          </motion.div>
        )}
      </div>
      <button
        className={styles.hamburger}
        aria-label="Open menu"
        aria-expanded={menuOpen}
        onClick={() => setMenuOpen(v => !v)}
      >
        <span />
        <span />
        <span />
      </button>
      <AnimatePresence>
        {menuOpen && (
          <motion.div
            className={styles.mobileMenuWrapper}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.18 }}
          >
            <motion.ul
              className={styles.mobileMenu}
              initial={{ y: -40, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              exit={{ y: -40, opacity: 0 }}
              transition={{ type: 'spring', stiffness: 80, damping: 14 }}
            >
              {navLinks.map(link => (
                <li key={link.label}>
                  <a href={link.href} onClick={() => setMenuOpen(false)}>{link.label}</a>
                </li>
              ))}
              <li className={styles.mobileActions}>
                {!isProfile && (
                  <>
                    <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
                      <Link to="/signup" className={styles.signUpBtn} tabIndex={0} onClick={() => setMenuOpen(false)}>
                        Sign Up
                      </Link>
                    </motion.div>
                    <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
                      <Link to="/login" className={styles.signInBtn} tabIndex={0} onClick={() => setMenuOpen(false)}>
                        Sign In
                      </Link>
                    </motion.div>
                  </>
                )}
                {isProfile && (
                  <motion.div whileHover={{ scale: 1.07 }} whileTap={{ scale: 0.97 }}>
                    <button className={styles.signInBtn} onClick={onLogout} tabIndex={0}>
                      Logout
                    </button>
                  </motion.div>
                )}
              </li>
            </motion.ul>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
};

export default Navbar; 