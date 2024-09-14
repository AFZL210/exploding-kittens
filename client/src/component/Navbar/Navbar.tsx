import React from 'react';
import './Navbar.css';
import { useNavigate, Link } from 'react-router-dom';

const Navbar: React.FC = () => {
  const navigate = useNavigate();

  return (
    <header className='header'>
      <nav>
        <h1 onClick={() => navigate('/')}>Exploding Kitten <img src='/kitten-logo.svg' className='kitten-logo'/></h1>
        <Link to='/leaderboard' className='nav-link'>Leaderboard</Link>
      </nav>
    </header>
  )
}

export default Navbar