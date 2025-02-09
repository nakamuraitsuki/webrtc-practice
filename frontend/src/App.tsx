import { BrowserRouter, Route, Routes } from 'react-router'
import { Home } from './pages/home'
import './App.css'

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home/>} />
          <Route path="/register" element={<div>Register</div>} />
          <Route path="/login" element={<div>Login</div>} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
