import reducer, {
  setAuth,
  logout
} from '@store/slices/auth'
import { AuthMode } from '@components/Auth/AuthCard'
import expect from 'expect'
// import { getAuth, connectAuthEmulator } from "firebase/auth";

// const auth = getAuth();
// connectAuthEmulator(auth, "http://localhost:9099");

const initialStateFixture = () => ({
  user: null,
  authMode: AuthMode.Login,
  initializing: false,
  vksp: ''
})

const userFixture = () => ({
  email: "test@email.com",
  uid: "1",
  name: "testName",
  claim: "testToken",
  photo: ""
})

describe('authSlice', function () {
  it('should return the default state', function () {
    const state = reducer(undefined, { type: "" })
    expect(state).toEqual(initialStateFixture())
  })

  it('should handle successful login', function () {
    const initialState = initialStateFixture()

    const expectedState = {
      ...initialState,
      user: userFixture(),
    }

    const state = reducer(initialState, setAuth(userFixture()))
    expect(state).toEqual(expectedState)
  })


  it('should handle logout', function () {
    const initialState = {
      ...initialStateFixture(),
      user: userFixture(),
    }
    const expected = initialStateFixture()
    // @ts-ignore
    const state = reducer(initialState, logout())
    expect(state).toEqual(expected)
  })
})