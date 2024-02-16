import '../App.css'

export const Footer = () => {


    return (
        <div className="footer">
            <a href="https://github.com/Anton-beep/distributed_calculations_go">
                <img src={process.env.PUBLIC_URL + '/github-mark.png'} alt="GitHub Mark"/>
            </a>
        </div>
    )
}