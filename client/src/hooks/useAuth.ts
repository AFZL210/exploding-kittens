import { useDispatch, useSelector } from "react-redux";
import { RootState, AppDispatch } from "../redux/store";
import { logoutUser, loginUser } from "../redux/slices/userSlice";

const useAuth = () => {
  const user = useSelector((state: RootState) => state.user);
  const dispatch = useDispatch<AppDispatch>();
  const isLoggedIn = user.isLoggedIn;
  const loading = user.loading;

  const logout = () => {
    dispatch(logoutUser());
  };

  const login = (username: string) => {
    dispatch(loginUser({ isLoggedIn: true, username: username, loading: true }));
  };

  return { isLoggedIn, login, logout, loading, user };
};

export default useAuth;
