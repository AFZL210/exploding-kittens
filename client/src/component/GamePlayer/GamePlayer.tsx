import React from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import { CardTypes } from "../../types/types";
import useAuth from "../../hooks/useAuth";

const GamePlayer: React.FC = () => {
  const { user } = useAuth();

  return (
    <div className="game-player">
      <div className="draw-count-container">
        <h4>{user.username}:</h4>
        <h4>Cards Drawn: 3/4</h4>
        <h4>Games Won: 4</h4>
        <h4>Games Lost: 3</h4>
      </div>

      <div className="cards-container">
        <Card cardType={CardTypes.Cat} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Bomb} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Defuse} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Shuffle} isFlipped={true} index={0} />
      </div>
    </div>
  );
};

export default GamePlayer;
