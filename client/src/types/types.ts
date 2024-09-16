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

export interface Card {
  cardType: CardTypes;
  isFlipped: boolean;
}

export interface GameStateI {
  cards: Card[];
  defuseCards: number;
  remainingCards: number;
  isLost: boolean;
  isWon: boolean;
  isLoading?: boolean;
}

export interface LeaderboardEntryI {
  username: string;
  score: number;
}

export interface openCardI {
  username: string;
  index: number;
}
