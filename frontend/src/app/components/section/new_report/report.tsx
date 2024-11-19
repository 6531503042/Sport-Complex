import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Newspaper,
  Calendar,
  Tag,
  ChevronRight,
  Trophy,
  Users,
  Star,
  TrendingUp,
} from "lucide-react";
import styles from './report.module.css';

interface NewsItem {
  id: string;
  date: string;
  title: string;
  description: string;
  category: string;
  image: string;
  author: string;
  readTime: string;
  tags: string[];
}

const newsData: NewsItem[] = [
  {
    id: "1",
    date: "2024-03-20",
    title: "New Fitness Center Equipment Arrival",
    description: "Experience our latest state-of-the-art gym equipment! We've added new treadmills, rowing machines, and a complete set of smart fitness equipment to enhance your workout experience.",
    category: "Facility Update",
    image: "https://images.unsplash.com/photo-1534438327276-14e5300c3a48",
    author: "Sport Complex Team",
    readTime: "3 min read",
    tags: ["Fitness", "Equipment", "Gym"]
  },
  {
    id: "2",
    date: "2024-03-18",
    title: "Swimming Competition Championship",
    description: "Join us for the annual MFU Swimming Championship! Open for all students and staff. Multiple categories available with exciting prizes to be won.",
    category: "Events",
    image: "https://images.unsplash.com/photo-1519315901367-f34ff9154487",
    author: "Events Team",
    readTime: "4 min read",
    tags: ["Swimming", "Competition", "Championship"]
  },
  {
    id: "3",
    date: "2024-03-15",
    title: "New Badminton Court Renovation",
    description: "We're excited to announce the completion of our badminton court renovation. New flooring, improved lighting, and enhanced ventilation systems have been installed.",
    category: "Facility Update",
    image: "https://images.unsplash.com/photo-1613918108466-292b78a8ef95",
    author: "Maintenance Team",
    readTime: "2 min read",
    tags: ["Badminton", "Renovation", "Facility"]
  },
  {
    id: "4",
    date: "2024-03-12",
    title: "Football Tournament Registration Open",
    description: "Register now for the upcoming inter-faculty football tournament. Form your team and compete for the championship title!",
    category: "Events",
    image: "https://images.unsplash.com/photo-1529900748604-07564a03e7a6",
    author: "Sports Committee",
    readTime: "5 min read",
    tags: ["Football", "Tournament", "Registration"]
  }
];

const Report: React.FC = () => {
  const [selectedNews, setSelectedNews] = useState<NewsItem>(newsData[0]);
  const [activeIndex, setActiveIndex] = useState<number>(0);

  const handleNewsClick = (news: NewsItem, index: number) => {
    setSelectedNews(news);
    setActiveIndex(index);
  };

  const categoryIcons = {
    "Facility Update": <Star className="w-5 h-5" />,
    "Events": <Trophy className="w-5 h-5" />,
    "Announcement": <Users className="w-5 h-5" />,
    "Achievement": <TrendingUp className="w-5 h-5" />
  };

  return (
    <div className={styles.reportContainer}>
      <motion.div className={styles.reportWrapper}>
        <div className={styles.reportGrid}>
          {/* News List Section */}
          <motion.div 
            className={styles.newsList}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
          >
            <h2 className={styles.newsListHeader}>Latest Updates</h2>
            <div className="space-y-4">
              {newsData.map((news, index) => (
                <motion.div
                  key={news.id}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  className={`cursor-pointer rounded-xl p-4 transition-all duration-200 ${
                    activeIndex === index 
                      ? 'bg-red-50 border-l-4 border-red-500 shadow-md' 
                      : 'hover:bg-gray-50 border-l-4 border-transparent'
                  }`}
                  onClick={() => handleNewsClick(news, index)}
                >
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-white shadow-sm">
                      {categoryIcons[news.category as keyof typeof categoryIcons]}
                    </div>
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900 mb-1">{news.title}</h3>
                      <div className="flex items-center gap-3 text-sm text-gray-500">
                        <span className="flex items-center gap-1">
                          <Calendar className="w-4 h-4" />
                          {new Date(news.date).toLocaleDateString()}
                        </span>
                        <span className="flex items-center gap-1">
                          <Tag className="w-4 h-4" />
                          {news.category}
                        </span>
                      </div>
                    </div>
                    <ChevronRight className={`w-5 h-5 transition-colors ${
                      activeIndex === index ? 'text-red-500' : 'text-gray-300'
                    }`} />
                  </div>
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* Selected News Detail */}
          <motion.div 
            className={styles.newsDetail}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
          >
            <AnimatePresence mode="wait">
              <motion.div
                key={selectedNews.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                className="bg-white rounded-2xl shadow-lg overflow-hidden"
              >
                <div className="relative h-64 bg-gray-200">
                  <div 
                    className="absolute inset-0 bg-cover bg-center"
                    style={{ backgroundImage: `url(${selectedNews.image})` }}
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                  <div className="absolute bottom-0 left-0 right-0 p-6 text-white">
                    <div className="flex items-center gap-2 mb-2">
                      <span className="px-3 py-1 rounded-full bg-white/20 backdrop-blur-sm text-sm">
                        {selectedNews.category}
                      </span>
                      <span className="text-sm">{selectedNews.readTime}</span>
                    </div>
                    <h1 className="text-3xl font-bold">{selectedNews.title}</h1>
                  </div>
                </div>

                <div className="p-6">
                  <div className="flex items-center gap-4 mb-6 text-sm text-gray-500">
                    <span className="flex items-center gap-2">
                      <Calendar className="w-4 h-4" />
                      {new Date(selectedNews.date).toLocaleDateString()}
                    </span>
                    <span>|</span>
                    <span>{selectedNews.author}</span>
                  </div>

                  <p className="text-gray-600 leading-relaxed mb-6">
                    {selectedNews.description}
                  </p>

                  <div className="flex flex-wrap gap-2">
                    {selectedNews.tags.map((tag) => (
                      <span 
                        key={tag}
                        className="px-3 py-1 rounded-full bg-gray-100 text-gray-600 text-sm"
                      >
                        #{tag}
                      </span>
                    ))}
                  </div>
                </div>
              </motion.div>
            </AnimatePresence>
          </motion.div>
        </div>
      </motion.div>
    </div>
  );
};

export default Report;
