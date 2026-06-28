import type React from "react";
import { Link } from "@tanstack/react-router";
import { sidebarMenu } from "@/constants/sideBarMenu";
import { CaretRight, CaretLeft, User } from "@phosphor-icons/react";

type SidebarProps = {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function Sidebar({ isOpen, setIsOpen }: SidebarProps) {
  return (
    <aside
      className={`fixed top-0 left-0 h-dvh bg-basic text-background text-white flex flex-col z-50 transition-all duration-300 ease-linear border-r ${
        isOpen ? "w-60" : "w-16"
      }`}
    >
      {/* Header with Logo */}
      <div className="p-4 border-b flex items-center gap-3">
        <div className="w-8 h-8 rounded flex items-center justify-center">
          <button
            className="block bg-basic hover:bg-primary rounded p-4"
            onClick={() => setIsOpen((prev) => !prev)}
          >
            {isOpen ? (
              <CaretLeft className="text-sm" />
            ) : (
              <CaretRight className="text-sm" />
            )}
          </button>
        </div>
        <div
          className={`transition-all ease-in duration-300 overflow-hidden ${
            isOpen ? "opacity-100 w-full" : "opacity-0 w-0"
          }`}
        >
          <div className="text-xl whitespace-nowrap">Tunnel Mesh</div>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 flex flex-col gap-1 p-2">
        {sidebarMenu.map((item) => {
          const IconComponent = item.icon;
          return (
            <Link
              key={item.label}
              to={item.to}
              className="group flex items-center gap-3 hover:bg-gray-700 p-3 hover:border-r rounded transition-colors relative hover:bg-primary hover:text-background"
            >
              <IconComponent className="text-xl flex-shrink-0" />
              <span
                className={`transition-all duration-300 whitespace-nowrap overflow-hidden ${
                  isOpen
                    ? "opacity-100 translate-x-0 w-auto"
                    : "opacity-0 -translate-x-2 w-0"
                }`}
              >
                {item.label}
              </span>
              {/* Tooltip for collapsed state */}
              {!isOpen && (
                <div className="absolute left-full ml-2 px-2 py-1 bg-gray-900 text-white text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                  {item.label}
                </div>
              )}
            </Link>
          );
        })}
      </nav>

      {/* Footer Profile */}
      <div className="p-4 border-t hover:bg-primary flex items-center gap-3">
        <div className="w-8 h-8 rounded flex items-center justify-center flex-shrink-0">
          <User />
        </div>
        <div
          className={`transition-all duration-300 overflow-hidden ${isOpen ? "opacity-100 w-full" : "opacity-0 w-0"}`}
        >
          <div className="text-sm whitespace-nowrap">Profile</div>
        </div>
      </div>
    </aside>
  );
}
