import React, { useState } from 'react';
import './Registration.css'
import Header from '../components/layout/Header/Header'

function Registration() {
    return(
        <div className="all_page">
            <Header />
            <RegistrationAndLogin />
        </div>
    );
}

function RegistrationAndLogin() {
    const [isLogIn, setIsLogIn] = useState(false);
    const [isActive, setIsActive] = useState(false);
    const [phone, setPhone] = useState('+7');

    const handleToggleForm = () => {
        setIsLogIn(!isLogIn);
    };

    const handleSubmitForm = () => {
        setIsActive(true);
    };

    // Обработчик изменения номера телефона
    const handlePhoneChange = (e) => {
        let value = e.target.value;
        
        if (!value.startsWith('+7')) {
            value = '+7';
        }
        
        // Удаляем все нецифровые символы, кроме +
        const digits = value.replace(/[^\d+]/g, '');
        
        if (digits.length <= 12) { // +7 + 10 цифр
            setPhone(digits);
        }
    };

    const handlePhoneFocus = (e) => {
        if (e.target.value === '') {
            setPhone('+7');
        }
    };

    return(
        <div className="container_login">  
            <div className="table">
                <div className="table-cell">
                    <div className={`container ${isLogIn ? 'log-in' : ''} ${isActive ? 'active' : ''}`}>
                        <div className="box"></div>
                        <div className="container-forms">
                            <div className="container-info">
                                <div className="info-item">
                                    <div className="table">
                                        <div className="table-cell">
                                            <p>Have an account?</p>
                                            <div className="btn" onClick={handleToggleForm}>Log in</div>
                                        </div>
                                    </div>
                                </div>
                                <div className="info-item">
                                    <div className="table">
                                        <div className="table-cell">
                                            <p>Don't have an account?</p>
                                            <div className="btn" onClick={handleToggleForm}>Sign up</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="container-form">
                                <div className="form-item log-in">
                                    <div className="table">
                                        <div className="table-cell">
                                            <input type="text" name="Username" placeholder="Логин" />
                                            <input type="password" name="Password" placeholder="Пароль" />
                                            <div className="btn" onClick={handleSubmitForm}>Log in</div>
                                        </div>
                                    </div>
                                </div>
                                <div className="form-item sign-up">
                                    <div className="table">
                                        <div className="table-cell">
                                            <input type="text" name="email" placeholder="Email" />
                                            <input 
                                                type="tel" 
                                                name="phone" 
                                                placeholder="+7 (XXX) XXX-XX-XX"
                                                value={phone}
                                                onChange={handlePhoneChange}
                                                onFocus={handlePhoneFocus}
                                                pattern="\+7\s?[\(]{0,1}[0-9]{3}[\)]{0,1}\s?\d{3}[-]{0,1}\d{2}[-]{0,1}\d{2}"
                                            />
                                            <input type="text" name="Username" placeholder="Логин" />
                                            <input type="password" name="Password" placeholder="Пароль" />
                                            <div className="btn" onClick={handleSubmitForm}>Sign up</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Registration