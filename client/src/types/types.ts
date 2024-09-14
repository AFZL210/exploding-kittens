export enum CardTypes {
  Cat = "Cat",
  Bomb = "Bomb",
  Defuse = "Defuse",
  Shuffle = "Shuffle",
}

export interface CardProps {
  isFlipped: boolean;
  cardType: CardTypes;
}
