import ky from "ky";

type BuildClientInput = {
  host?: string;
  cookie?: string;
};

const buildClient = (clientParams?: BuildClientInput) => {
  if (typeof window === "undefined") {
    if (!clientParams) throw new Error("Client params not provided");
    const { host, cookie } = clientParams;
    // On server, we manually specify host and cookie headers
    return ky.create({
      prefixUrl:
        "http://ingress-nginx-controller.ingress-nginx.svc.cluster.local",
      headers: {
        Host: host,
        Cookie: cookie,
      },
    });
  } else {
    // Browser takes care of handling host and cookie for us
    return ky.create({
      prefixUrl: "/",
    });
  }
};

export { buildClient };
