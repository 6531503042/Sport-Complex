import React from "react";
import Link from "next/link";
import Logo from "../../../assets/transparent.png";
import { motion } from "framer-motion";
import { 
  Facebook, 
  Instagram, 
  Twitter, 
  YouTube,
  Email,
  Phone,
  LocationOn
} from "@mui/icons-material";

const Footer = () => {
  const socialLinks = [
    { icon: <Facebook />, href: "https://facebook.com" },
    { icon: <Instagram />, href: "https://instagram.com" },
    { icon: <Twitter />, href: "https://twitter.com" },
    { icon: <YouTube />, href: "https://youtube.com" }
  ];

  const contactInfo = [
    { icon: <Email />, text: "contact@sportcomplex.com" },
    { icon: <Phone />, text: "+66 123 456 789" },
    { icon: <LocationOn />, text: "Mae Fah Luang University, Chiang Rai" }
  ];

  return (
    <motion.footer 
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.8 }}
      className="bg-gradient-to-b from-slate-800 to-slate-900 text-white"
    >
      <div className="max-w-7xl mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
          {/* Logo and Description */}
          <motion.div 
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.2 }}
            className="space-y-4"
          >
            <Link href="/" className="inline-flex items-center gap-3.5">
              <img src={Logo.src} alt="Logo" className="w-10 h-10" />
              <div className="border-l-2 border-white/20 pl-3">
                <div className="font-semibold text-xl">
                  <span className="text-white">SPORT.</span>
                  <span className="text-gray-400">MFU</span>
                </div>
                <div className="text-sm text-gray-400">SPORT COMPLEX</div>
              </div>
            </Link>
            <p className="text-gray-400 mt-4 max-w-md">
              Your premier destination for sports and fitness at Mae Fah Luang University. 
              Experience top-notch facilities and professional services.
            </p>
          </motion.div>

          {/* Quick Links */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
            className="space-y-4"
          >
            <h3 className="text-xl font-semibold mb-4">Quick Links</h3>
            <ul className="space-y-2">
              {['About Us', 'Facilities', 'Membership', 'Events', 'Contact'].map((item) => (
                <motion.li 
                  key={item}
                  whileHover={{ x: 5 }}
                  className="text-gray-400 hover:text-white transition-colors"
                >
                  <Link href="/">{item}</Link>
                </motion.li>
              ))}
            </ul>
          </motion.div>

          {/* Contact Info */}
          <motion.div 
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.6 }}
            className="space-y-4"
          >
            <h3 className="text-xl font-semibold mb-4">Contact Us</h3>
            <div className="space-y-3">
              {contactInfo.map((item, index) => (
                <motion.div 
                  key={index}
                  whileHover={{ x: 5 }}
                  className="flex items-center gap-3 text-gray-400 hover:text-white transition-colors"
                >
                  {item.icon}
                  <span>{item.text}</span>
                </motion.div>
              ))}
            </div>

            {/* Social Links */}
            <div className="flex gap-4 mt-6">
              {socialLinks.map((social, index) => (
                <motion.a
                  key={index}
                  href={social.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  whileHover={{ scale: 1.2, rotate: 5 }}
                  whileTap={{ scale: 0.9 }}
                  className="bg-white/10 p-2 rounded-full hover:bg-white/20 transition-colors"
                >
                  {social.icon}
                </motion.a>
              ))}
            </div>
          </motion.div>
        </div>

        {/* Copyright */}
        <motion.div 
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ delay: 0.8 }}
          className="border-t border-white/10 mt-12 pt-8 text-center text-gray-400"
        >
          <p>© {new Date().getFullYear()} MFU Sport Complex. All rights reserved.</p>
        </motion.div>
      </div>
    </motion.footer>
  );
};

export default Footer;
