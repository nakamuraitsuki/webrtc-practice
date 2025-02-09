import { BrowserRouter, Route, Routes } from 'react-router'
import { Home } from './pages/home'
import { Register } from './pages/register'
import './App.css'

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home/>} />
          <Route path="/register" element={<Register/>} />
          <Route path="/login" element={<div>Login</div>} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
