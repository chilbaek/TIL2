import PropTypes from "prop-types"
import styeles from "./Button.module.css"

function Button({text}) {
    return (<button className={styeles.btn}>{text}</button>)
}

Button.propTypes = {
    text: PropTypes.string.isRequired,
};

export default Button;
