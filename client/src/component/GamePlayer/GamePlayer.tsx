import React, { useEffect, useState } from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import useAuth from "../../hooks/useAuth";
import { useDispatch, useSelector } from "react-redux";
import { fetchGameState, openCard } from "../../redux/slices/gameSlice";
import { RootState } from "../../redux/store";
import toast from "react-hot-toast";

const GamePlayer: React.FC = () => {
  const { user } = useAuth();
  const dispatch = useDispatch();
  const gameState = useSelector((state: RootState) => state.game);

  const [isDisabled, setIsDisabled] = useState(false);
  const [toastShown, setToastShown] = useState(false);

  const getGameStateDataAsync = () => {
    // @ts-ignore
    dispatch(fetchGameState(user.username));
  };

  const openCardAsync = (index: number) => {
    if (isDisabled) return;
    // @ts-ignore
    dispatch(openCard({ index: index, username: user.username }));
  };

  useEffect(() => {
    getGameStateDataAsync();
  }, []);

  useEffect(() => {
    if ((gameState.isLost || gameState.isWon) && !toastShown) {
      setIsDisabled(true);
      gameState.isLost ? toast.error("You Lost!") : toast.success("You Won!");
      setToastShown(true);

      setTimeout(() => {
        setIsDisabled(false);
        setToastShown(false);
        getGameStateDataAsync();
      }, 2000);
    }
  }, [gameState, toastShown]);

  return (
    <div className={`game-player ${isDisabled ? "fade-out" : ""}`}>
      <div className="draw-count-container">
        <h4>{user.username}:</h4>
        <h4>Cards Drawn: {5 - gameState.remainingCards}/5</h4>
      </div>

      <div className="cards-container">
        {gameState.cards.map((card, idx: number) => (
          <div
            onClick={() =>
              !gameState.cards[idx].isFlipped && !isDisabled && openCardAsync(idx)
            }
            key={idx}
          >
            <Card
              cardType={card.cardType}
              index={idx}
              isFlipped={card.isFlipped}
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default GamePlayer;
