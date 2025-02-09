import { useState } from 'react';

export const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setMessage('');
        setLoading(true);

        try {
            const response = await fetch('http://localhost:8080/api/user/authenticate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`ログイン失敗: ${errorText}`);
            }

            setMessage("ログインに成功しました");
        } catch (error) {
            if (error instanceof TypeError) {
                setMessage('error:' + error.message);
            } else {
                setMessage('unknown error' + error);
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <input
                    type="mail"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button type="submit" disabled={loading} style={{ padding: "10px 20px" }}>
                    {loading ? "ログイン..." : "ログイン"}
                </button>
            </form>
            {message && <p style={{ marginTop: "10px", color: "red" }}>{message}</p>}
        </div>
    )
}