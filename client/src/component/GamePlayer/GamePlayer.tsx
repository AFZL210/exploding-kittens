import React, { useEffect, useState } from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import useAuth from "../../hooks/useAuth";
import { useDispatch, useSelector } from "react-redux";
import { fetchGameState, openCard } from "../../redux/slices/gameSlice";
import { RootState } from "../../redux/store";

const GamePlayer: React.FC = () => {
  const { user } = useAuth();
  const dispatch = useDispatch();
  const gameState = useSelector((state: RootState) => state.game);

  const getGameStateDataAsync = () => {
    // @ts-ignore
    dispatch(fetchGameState(user.username));
  };

  const openCardAsync = (index: number) => {
    // @ts-ignore
    dispatch(openCard({ index: index, username: user.username }));
  };

  useEffect(() => {
    getGameStateDataAsync();
  }, []);

  useEffect(() => {
    if(gameState.isLost) {
      alert('lost')
    }
    if(gameState.isWon) {
      alert('won')
    }
  }, [gameState]);

  return (
    <div className={`game-player`}>
      <div className="draw-count-container">
        <h4>{user.username}:</h4>
        <h4>Cards Drawn: {5 - gameState.remainingCards}/5</h4>
      </div>

      <div className="cards-container">
        {gameState.cards.map((card, idx: number) => {
          return (
            <div
              onClick={() =>
                !gameState.cards[idx].isFlipped && openCardAsync(idx)
              }
              key={idx}
            >
              <Card
                cardType={card.cardType}
                index={idx}
                isFlipped={card.isFlipped}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default GamePlayer;
