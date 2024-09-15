import React, { FormEvent, useState } from "react";
import "./Login.css";
import toast from "react-hot-toast";
import useAuth from "../../hooks/useAuth";

const Login: React.FC = () => {
  const [showLoginTab, setShowLoginTab] = useState(true);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { loading } = useAuth();

  const loginUser = async (e: FormEvent) => {
    e.preventDefault();
    toast("hello");
  };

  const registerUser = async (e: FormEvent) => {
    e.preventDefault();
  };

  const toggleTab = () => {
    setUsername("");
    setPassword("");
    setShowLoginTab(!showLoginTab);
  };

  return (
    <div className="login-container">
      {showLoginTab ? (
        <form className="form-container" onSubmit={loginUser}>
          <h1>Login</h1>
          <div className="input-control">
            <input type="text" placeholder="username" />
          </div>
          <div className="input-control">
            <input type="password" placeholder="password" />
          </div>
          <button type="submit" className="submit-btn" disabled={loading}>
            Login
          </button>
          <p>
            Dont't have an account?{" "}
            <span className="toggle-link" onClick={() => toggleTab()}>
              Register
            </span>
          </p>
        </form>
      ) : (
        <form className="form-container" onSubmit={registerUser}>
          <h1>Register</h1>
          <div className="input-control">
            <input type="text" placeholder="username" />
          </div>
          <div className="input-control">
            <input type="password" placeholder="password" />
          </div>
          <button type="submit" className="submit-btn" disabled={loading}>
            Register
          </button>
          <p>
            Already have an account?{" "}
            <span className="toggle-link" onClick={() => toggleTab()}>
              Login
            </span>
          </p>
        </form>
      )}
    </div>
  );
};

export default Login;
