
import sinon = require("sinon");
import BluebirdPromise = require("bluebird");
import Assert = require("assert");
import U2FRegisterRequestGet = require("../../../../../src/lib/routes/secondfactor/u2f/register_request/get");
import AuthenticationSession = require("../../../../../src/lib/AuthenticationSession");
import { ServerVariablesHandler } from "../../../../../src/lib/ServerVariablesHandler";
import winston = require("winston");

import ExpressMock = require("../../../../mocks/express");
import { UserDataStoreStub } from "../../../../mocks/storage/UserDataStoreStub";
import U2FMock = require("../../../../mocks/u2f");
import ServerVariablesMock = require("../../../../mocks/ServerVariablesMock");
import U2f = require("u2f");

describe("test u2f routes: register_request", function () {
  let req: ExpressMock.RequestMock;
  let res: ExpressMock.ResponseMock;
  let mocks: ServerVariablesMock.ServerVariablesMock;
  let authSession: AuthenticationSession.AuthenticationSession;

  beforeEach(function () {
    req = ExpressMock.RequestMock();
    req.app = {};
    mocks = ServerVariablesMock.mock(req.app);
    req.session = {};
    AuthenticationSession.reset(req as any);

    req.headers = {};
    req.headers.host = "localhost";

    const options = {
      inMemoryOnly: true
    };


    mocks.userDataStore.saveU2FRegistrationStub.returns(BluebirdPromise.resolve({}));
    mocks.userDataStore.retrieveU2FRegistrationStub.returns(BluebirdPromise.resolve({}));

    res = ExpressMock.ResponseMock();
    res.send = sinon.spy();
    res.json = sinon.spy();
    res.status = sinon.spy();

    return AuthenticationSession.get(req as any)
      .then(function (_authSession: AuthenticationSession.AuthenticationSession) {
        authSession = _authSession;
        authSession.userid = "user";
        authSession.first_factor = true;
        authSession.second_factor = false;
        authSession.identity_check = {
          challenge: "u2f-register",
          userid: "user"
        };
      });
  });

  describe("test registration request", () => {
    it("should send back the registration request and save it in the session", function () {
      const expectedRequest = {
        test: "abc"
      };
      const user_key_container = {};
      const u2f_mock = U2FMock.U2FMock();
      u2f_mock.request.returns(BluebirdPromise.resolve(expectedRequest));

      mocks.u2f = u2f_mock;
      return U2FRegisterRequestGet.default(req as any, res as any)
        .then(function () {
          Assert.deepEqual(expectedRequest, res.json.getCall(0).args[0]);
        });
    });

    it("should return internal error on registration request", function () {
      res.send = sinon.spy();
      const user_key_container = {};
      const u2f_mock = U2FMock.U2FMock();
      u2f_mock.request.returns(BluebirdPromise.reject("Internal error"));

      mocks.u2f = u2f_mock;
      return U2FRegisterRequestGet.default(req as any, res as any)
      .then(function() {
        Assert.equal(res.status.getCall(0).args[0], 200);
        Assert.deepEqual(res.send.getCall(0).args[0], {
          error: "Operation failed."
        });
      });
    });

    it("should return forbidden if identity has not been verified", function () {
      authSession.identity_check = undefined;
      return U2FRegisterRequestGet.default(req as any, res as any)
        .then(function () {
          Assert.equal(403, res.status.getCall(0).args[0]);
        });
    });
  });
});
