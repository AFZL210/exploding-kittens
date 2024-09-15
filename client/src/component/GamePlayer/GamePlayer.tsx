import React, { useEffect, useState } from "react";
import "./GamePlayer.css";
import Card from "../Card/Card";
import useAuth from "../../hooks/useAuth";
import { useDispatch, useSelector } from "react-redux";
import { fetchGameState, openCard } from "../../redux/slices/gameSlice";
import { RootState } from "../../redux/store";
import { CardTypes } from "../../types/types";
import toast from "react-hot-toast";

const GamePlayer: React.FC = () => {
  const { user } = useAuth();
  const dispatch = useDispatch();
  const gameState = useSelector((state: RootState) => state.game);

  const [isDisabled, setIsDisabled] = useState(false);
  const [fadeClass, setFadeClass] = useState("");

  const getGameStateDataAsync = () => {
    // @ts-ignore
    dispatch(fetchGameState(user.username));
  };

  const openCardAsync = (index: number) => {
    if (isDisabled) return;

    setIsDisabled(true);
    const cardType = gameState.cards[index].cardType;

    // @ts-ignore
    dispatch(openCard({ index: index, username: user.username }));
    if (
      [CardTypes.Shuffle, CardTypes.Bomb].includes(
        cardType
      )
    ) {
      toast(`It's a ${gameState.cards[index].cardType} Card. You lost!`);
      setFadeClass("fade-out");

      new Promise((res, _rej) => {
        setTimeout(() => {
          res(null);
        }, 5000);
      }).then(() => {
        setFadeClass("");
        setIsDisabled(false);
        getGameStateDataAsync();
      });
    } else {
      setIsDisabled(false);
    }
  };

  useEffect(() => {
    getGameStateDataAsync();
  }, []);

  useEffect(() => {
    console.log(gameState);
  }, [gameState]);

  return (
    <div className={`game-player ${fadeClass}`}>
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
              style={{ pointerEvents: isDisabled ? "none" : "auto" }}
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
