import { useDispatch, useSelector } from "react-redux";
import { RootState, AppDispatch } from "../redux/store";
import { loginUser } from "../redux/slices/userSlice";

const useAuth = () => {
  const user = useSelector((state: RootState) => state.user);
  const dispatch = useDispatch<AppDispatch>();
  const isLoggedIn = user.isLoggedIn;
  const loading = user.loading;

  const login = (username: string) => {
    dispatch(
      loginUser({ isLoggedIn: true, username: username, loading: true })
    );
    localStorage.setItem("username", username);
  };
  return { isLoggedIn, login, loading, user };
};

export default useAuth;
