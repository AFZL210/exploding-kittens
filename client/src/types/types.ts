export enum CardTypes {
  Cat = "Cat",
  Bomb = "Bomb",
  Defuse = "Defuse",
  Shuffle = "Shuffle",
}

export interface CardI {
  index: number;
  isFlipped: boolean;
  cardType: CardTypes;
}

export interface UserI {
  isLoggedIn: boolean;
  username: string;
  loading: boolean;
}
