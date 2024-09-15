import React, { useEffect } from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import { CardTypes } from "../../types/types";
import useAuth from "../../hooks/useAuth";
import { useDispatch } from "react-redux";
import { fetchGameState } from "../../redux/slices/gameSlice";

const GamePlayer: React.FC = () => {
  const { user } = useAuth();
  const dispatch = useDispatch();
  
  const getGameState = async () => {
    dispatch(fetchGameState())
    console.log("hello")
  }

  useEffect(() => {
    getGameState();
  }, [])

  return (
    <div className="game-player">
      <div className="draw-count-container">
        <h4>{user.username}:</h4>
        <h4>Cards Drawn: 3/4</h4>
        <h4>Games Won: 4</h4>
        <h4>Games Lost: 3</h4>
      </div>

      <div className="cards-container">
        <button onClick={() => getGameState()}>click</button>
        <Card cardType={CardTypes.Cat} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Bomb} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Defuse} isFlipped={true} index={0} />
        <Card cardType={CardTypes.Shuffle} isFlipped={true} index={0} />
      </div>
    </div>
  );
};

export default GamePlayer;
