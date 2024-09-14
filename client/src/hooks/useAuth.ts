import { useDispatch, useSelector } from "react-redux";
import { RootState, AppDispatch } from "../redux/store";
import { logoutUser, loginUser } from "../redux/slices/userSlice";

const useAuth = () => {
  const user = useSelector((state: RootState) => state.user);
  const dispatch = useDispatch<AppDispatch>();
  const isLoggedIn = user.isLoggedIn;

  const logout = () => {
    dispatch(logoutUser());
  };

  const login = (username: string) => {
    dispatch(loginUser({ isLoggedIn: true, username: username }));
  };

  return { isLoggedIn, login, logout };
};

export default useAuth;
