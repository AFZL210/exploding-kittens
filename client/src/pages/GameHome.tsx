import React from "react";
import GamePlayer from "../component/GamePlayer/GamePlayer";
import Login from "../component/Login/Login";
import useAuth from "../hooks/useAuth";

const GameHome: React.FC = () => {
  const {isLoggedIn} = useAuth();
  
  return (
    <div className="home-container">
      {isLoggedIn ? <GamePlayer /> : <Login />}
    </div>
  );
};

export default GameHome;
