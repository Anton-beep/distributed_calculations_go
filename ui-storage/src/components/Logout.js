import Cookies from "js-cookie";

export const Logout = () => {
    Cookies.remove('token');
    window.location.reload()

    return (
        <div className="container">
            <h1>Logout</h1>
            <p>You have been logged out</p>
        </div>
    );
}