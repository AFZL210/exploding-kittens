import React, { FormEvent, useState } from "react";
import "./Login.css";
import toast from "react-hot-toast";
import useAuth from "../../hooks/useAuth";
import apiClient from "../../lib/apiClient";

const Login: React.FC = () => {
  const [showLoginTab, setShowLoginTab] = useState(true);
  const [errorMessage, setErrorMessage] = useState<null | string>(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { loading, login } = useAuth();

  const resetUserData = () => {
    setUsername("");
    setPassword("");
  };

  const validateUserData = () => {
    if (username === "") {
      setErrorMessage("Username cannot be empty.");
      return false;
    } else if (password === "") {
      setErrorMessage("Password cannot be empty.");
      return false;
    }

    return true;
  };

  const loginUser = async (e: FormEvent) => {
    e.preventDefault();
    if (!validateUserData()) return;
    setErrorMessage(null);

    apiClient
      .post("/login", { username, password })
      .then(() => {
        toast("Logged in!");
        login(username);
        localStorage.setItem("username", username);
        resetUserData();
      })
      .catch(() => {
        toast("Wrong password or username");
      });
  };

  const registerUser = async (e: FormEvent) => {
    e.preventDefault();
    if (!validateUserData()) return;
    setErrorMessage(null);

    apiClient
      .post("/register", {
        username,
        password,
      })
      .then(() => {
        toast("Registered!");
        resetUserData();
        setShowLoginTab(true);
      })
      .catch(() => {
        toast("User already exists");
      });
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
          {errorMessage !== null && <span>{errorMessage}</span>}
          <div className="input-control">
            <input
              type="text"
              placeholder="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className="input-control">
            <input
              type="password"
              placeholder="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
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
          {errorMessage !== null && <span>{errorMessage}</span>}
          <div className="input-control">
            <input
              type="text"
              placeholder="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className="input-control">
            <input
              type="password"
              placeholder="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
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
