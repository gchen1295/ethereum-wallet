// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";

// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
export const firebaseConfig = {
  apiKey: "AIzaSyBvN37u78QbbrXGGW42e5algszc8JdaY0Y",
  authDomain: "nft-app-3db02.firebaseapp.com",
  projectId: "nft-app-3db02",
  storageBucket: "nft-app-3db02.appspot.com",
  messagingSenderId: "384145460687",
  appId: "1:384145460687:web:32b4af443bade0f11656ba",
  measurementId: "G-69X7L5S0XP"
};

// Initialize Firebase
export const app = initializeApp(firebaseConfig);
export const analytics = getAnalytics(app);