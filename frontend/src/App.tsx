import { BrowserRouter, Route, Routes } from 'react-router'
import { Home } from './pages/home'
import { Signup } from './pages/signup'
import { Login } from './pages/login'
import './App.css'

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home/>} />
          <Route path="/register" element={<Signup/>} />
          <Route path="/login" element={<Login/>} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
