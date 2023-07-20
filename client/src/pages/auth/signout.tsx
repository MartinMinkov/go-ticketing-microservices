import { useEffect } from "react";
import Router from "next/router";
import useRequest from "@/hooks/use-request";

const SignOut = () => {
  const { doRequest } = useRequest({
    method: "post",
    url: "/api/users/signout",
  });

  useEffect(() => {
    async function fetchData() {
      await doRequest({}, undefined);
      Router.push("/");
    }
    fetchData();
  });

  return <div>Signing you out...</div>;
};

export default SignOut;
