/*
Author: KleaSCM
Email: KleaSCM@gmail.com
File: Footer.tsx
Description: Modern, animated footer for SteriaFront. Uses Framer Motion and modular SCSS for style and animation.
*/

import React from 'react';
import { motion } from 'framer-motion';
import styles from './Footer.module.scss';

const Footer: React.FC = () => (
  <motion.footer
    className={styles.footer}
    initial={{ y: 60, opacity: 0 }}
    animate={{ y: 0, opacity: 1 }}
    transition={{ type: 'spring', stiffness: 60, damping: 12, delay: 0.2 }}
  >
    <span className={styles.brand}>SteriaFront</span>
    <span className={styles.copyright}>
      &copy; {new Date().getFullYear()} All rights reserved.
    </span>
  </motion.footer>
);

export default Footer; 