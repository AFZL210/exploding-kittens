import UserTile from '../component/UserTile/UserTile';
import './LeaderBoard.css';

const LeaderBoard = () => {
    return (
        <div className="leaderboard">
            <h1>Leaderboard</h1>
            <h3>Your Rank: 10</h3>
            <UserTile />
        </div>
    )
}

export default LeaderBoard;