import { BrowserRouter, Route, Routes } from 'react-router'
import './App.css'

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<div>Home</div>} />
          <Route path="/login" element={<div>Login</div>} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
