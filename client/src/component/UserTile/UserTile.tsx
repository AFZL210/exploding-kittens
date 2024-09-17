import React, { useEffect, useState } from "react";
import "./UserTile.css";
import apiClient from "../../lib/apiClient";

interface UserRankI {
  username: string;
  score: number;
}

const UserTile: React.FC = () => {
  const [leaderboard, setLeaderboard] = useState<UserRankI[]>([]);
  const [loading, setLoading] = useState(true);

  const getLeaderBoardAsync = async () => {
    apiClient.get("/leaderboard").then((data) => {
      setLeaderboard(data.data);
      setLoading(false);
    });
  };

  useEffect(() => {
    getLeaderBoardAsync();

    const interval = setInterval(() => {
      getLeaderBoardAsync();
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="user-tile">
      <div className="rank-heading tile">
        <h4>Username</h4>
        <h4>Rank</h4>
      </div>
      {loading ? (
        <span>Loading...</span>
      ) : (
        leaderboard.map((userRank: UserRankI, idx: number) => {
          return (
            <div className="user-rank tile" key={idx}>
              <h4>{userRank.username}</h4>
              <h4>{idx+1}</h4>
            </div>
          );
        })
      )}
    </div>
  );
};

export default UserTile;