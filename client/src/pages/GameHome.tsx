import React from "react";
import GamePlayer from "../component/GamePlayer/GamePlayer";
import Login from "../component/Login/Login";

const GameHome: React.FC = () => {
  let isLoggedIn = true;
  return (
    <div className="home-container">
      {isLoggedIn ? <GamePlayer /> : <Login />}
    </div>
  );
};

export default GameHome;
