import { useRef, useState, useEffect, useReducer } from "react";
import { useFuzzySearchList } from "@nozbe/microfuzz/react";
import { FaUniversity } from "react-icons/fa";

import { GetLiterature, OpenPDF, Exit } from "../wailsjs/go/main/App";

import List from "./Components/List";
import { ItemProps } from "./Components/Item";

import "./App.css";

const App = () => {
  const [allitems, setAllItems] = useState<ItemProps[]>([]);
  const [debug, setDebug] = useState<string>("");
  const [prompt, setPrompt] = useState<string>("");
  const inputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    GetLiterature().then((res) => {
      setAllItems(res as unknown as ItemProps[]);
    });

    const handleBlur = () => {
      setTimeout(() => {
        inputRef.current?.focus({ preventScroll: true });
      }, 20);
    };

    const inputElement = inputRef.current;
    inputElement?.addEventListener("blur", handleBlur);

    return () => {
      inputElement?.removeEventListener("blur", handleBlur);
    };
  }, []);

  const filteredList = useFuzzySearchList({
    list: allitems,
    queryText: prompt,
    getText: (item) => [item.year, item.author, item.title],
    mapResultItem: ({ item, score, matches: [hlyear, hlauthor, hltitle] }) => {
      return {
        item,
        score,
        hl: {
          year: hlyear,
          author: hlauthor,
          title: hltitle,
        },
      };
    },
  });

  const [state, dispatch] = useReducer(
    (state: { active: number }, action: { type: string; payload?: any }) => {
      if (action.type === "arrowUp") {
        return {
          active:
            state.active <= 0 ? filteredList.length - 1 : state.active - 1,
        };
      } else if (action.type === "arrowDown") {
        return {
          active:
            state.active >= filteredList.length - 1 ? 0 : state.active + 1,
        };
      } else if (action.type === "escape") {
        Exit();
      } else if (action.type === "select") {
        OpenPDF(filteredList[state.active].item.file);
      } else if (action.type === "filter") {
        setPrompt(action.payload.target.value);
      }
      return state;
    },
    {
      active: 0,
    },
  );

  useEffect(() => {
    let intervalId: ReturnType<typeof setInterval> | null = null;
    let timeoutId: ReturnType<typeof setTimeout> | null = null;

    const startRepeatingAction = (key: string) => {
      const actionType = key === "ArrowUp" ? "arrowUp" : "arrowDown";
      dispatch({ type: actionType });

      timeoutId = setTimeout(() => {
        intervalId = setInterval(() => {
          dispatch({ type: actionType });
        }, 50);
      }, 300);
    };

    const stopRepeatingAction = () => {
      if (timeoutId) clearTimeout(timeoutId);
      if (intervalId) clearInterval(intervalId);
      timeoutId = null;
      intervalId = null;
    };

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "ArrowUp" || e.key === "ArrowDown") {
        e.preventDefault();
        if (!timeoutId && !intervalId) {
          startRepeatingAction(e.key);
        }
      } else if (e.key === "Escape") {
        dispatch({ type: "escape" });
      } else if (e.key === "Enter") {
        dispatch({ type: "select" });
      }
    };

    const handleKeyUp = (e: KeyboardEvent) => {
      if (e.key === "ArrowUp" || e.key === "ArrowDown") {
        stopRepeatingAction();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    window.addEventListener("keyup", handleKeyUp);

    return () => {
      stopRepeatingAction();
      window.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("keyup", handleKeyUp);
    };
  }, [dispatch]);

  return (
    <div id="App">
      <div className="navbar">
        <div className="debug">
          DebugINFO: {state.active} {debug}
        </div>
        <FaUniversity />
        <div id="input" className="input-box">
          <input
            placeholder="refs"
            id="prompt"
            className="input"
            onChange={(e) => {
              dispatch({ type: "filter", payload: e });
            }}
            autoComplete="off"
            name="input"
            type="text"
            ref={inputRef}
            autoFocus
          />
        </div>
      </div>
      <List itemlist={filteredList} active={state.active} />
    </div>
  );
};

export default App;
