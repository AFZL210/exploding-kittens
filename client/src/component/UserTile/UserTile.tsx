import React from 'react'
import './UserTile.css';

const UserTile: React.FC = () => {
  return (
    <div className='user-tile'>
        <div className='rank-heading tile'>
            <h4>Username</h4>
            <h4>Rank</h4>
        </div>
        <div className='user-rank tile'>
            <h4>Username</h4>
            <h4>Rank</h4>
        </div>
        <div className='user-rank tile'>
            <h4>Username</h4>
            <h4>Rank</h4>
        </div>
        <div className='user-rank tile'>
            <h4>Username</h4>
            <h4>Rank</h4>
        </div>
    </div>
  )
}

export default UserTile