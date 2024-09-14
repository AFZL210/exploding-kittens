import React from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import { CardTypes } from "../../types/types";

const GamePlayer: React.FC = () => {
  return (
    <div className="game-player">
      <div className="draw-count-container">
        <h4>AFZL210:</h4>
        <h4>Cards Drawn: 3/4</h4>
        <h4>Games Won: 4</h4>
        <h4>Games Lost: 3</h4>
      </div>

      <div className="cards-container">
        <Card cardType={CardTypes.Cat} isFlipped={true} />
        <Card cardType={CardTypes.Bomb} isFlipped={true} />
        <Card cardType={CardTypes.Defuse} isFlipped={true} />
        <Card cardType={CardTypes.Shuffle} isFlipped={true} />
      </div>
    </div>
  );
};

export default GamePlayer;
