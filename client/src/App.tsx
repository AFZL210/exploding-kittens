import React from "react";
import './App.css';
import { Route, Routes } from "react-router-dom";
import GameHome from "./pages/GameHome";
import LeaderBoard from "./pages/LeaderBoard";
import Navbar from "./component/Navbar/Navbar";

const App: React.FC = () => {
  return (
    <div className="app-container">
      <Navbar />
      <div className="pages">
        <Routes>
          <Route path="/" element={<GameHome />} />
          <Route path="/leaderboard" element={<LeaderBoard />} />
        </Routes>
      </div>
    </div>
  );
};

export default App;
