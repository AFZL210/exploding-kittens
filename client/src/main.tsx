import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from './redux/store.ts'
import { Toaster } from 'react-hot-toast'

createRoot(document.getElementById('root')!).render(
    <BrowserRouter>
      <Provider store={store}>
        <Toaster></Toaster>
        <App />
      </Provider>
    </BrowserRouter>
)


