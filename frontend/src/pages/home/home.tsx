import { Link } from 'react-router'

export const Home = () => {
    return (
        <div>
            <h1>Sample</h1>
            <ul>
                <li><Link to="/register">Register</Link></li>
                <li><Link to="/login">Login</Link></li>
            </ul>
        </div>
    )
}