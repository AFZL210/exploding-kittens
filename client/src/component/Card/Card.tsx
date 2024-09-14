import React from "react";
import "./Card.css";
import { CardI } from "../../types/types";

const imageTypePathMap: Record<string, string> = {
    Cat: '/kitten-logo.svg',
    Bomb: '/bomb-icon.svg',
    Defuse: '/defuse-icon.svg',
    Shuffle: '/shuffle-icon.png'
}

const Card: React.FC<CardI> = (props) => {
    return (
    <div className="card">
        {props.isFlipped ?
        <div className="flipped-side">
            <img src={imageTypePathMap[props.cardType]} alt="" />
            <h1>{props.cardType}</h1>
        </div>: 
        <div className="hidden-card">
            <h1>?</h1>
        </div>
        }
    </div>
  );
};

export default Card;
