import { useEffect, useState } from "react";
import UserTile from "../component/UserTile/UserTile";
import "./LeaderBoard.css";
import useAuth from "../hooks/useAuth";
import apiClient from "../lib/apiClient";

interface UserRank {
  rank: number;
  score: number;
  username: string;
}

const LeaderBoard = () => {
  const [userrank, setUserRank] = useState<null | number>(null);
  const { user } = useAuth();

  const getUserRank = async () => {
    if (user.isLoggedIn) {
      apiClient.get(`/user-rank?username=${user.username}`).then((data) => {
        setUserRank((data.data as UserRank).rank);
      });
    }
  };

  useEffect(() => {
    getUserRank();
  }, [])

  return (
    <div className="leaderboard">
      <h1>Leaderboard</h1>
      <h3>Your Rank: {userrank ?? "N/A"}</h3>
      <UserTile />
    </div>
  );
};

export default LeaderBoard;
