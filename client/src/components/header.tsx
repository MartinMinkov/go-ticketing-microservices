import { User } from "@/models/user";
import Link from "next/link";

type HeaderProps = {
  user: User | {};
};

const Header = ({ user }: HeaderProps) => {
  const renderLinks = (user: User | {}) => {
    if (Object.keys(user).length) {
      return (
        <li className="hover:text-blue-500">
          <Link href="/auth/signout">Sign Out</Link>
        </li>
      );
    } else {
      return (
        <>
          <li className="mr-6 hover:text-blue-500">
            <Link href="/auth/signup">Sign Up</Link>
          </li>
          <li className="hover:text-blue-500">
            <Link href="/auth/signin">Sign In</Link>
          </li>
        </>
      );
    }
  };

  return (
    <nav className="rounded border-gray-200 bg-gray-200 px-2 py-2.5 sm:px-4">
      <div className="container flex flex-wrap items-center justify-between">
        <Link href="/">GitTix</Link>
        <div>
          <ul className="flex items-center ">{renderLinks(user)}</ul>
        </div>
      </div>
    </nav>
  );
};

export default Header;
