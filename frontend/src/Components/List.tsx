import { HighlightRanges } from "@nozbe/microfuzz";

import Item from "./Item";
import { ItemProps } from "./Item";
import React, { useRef, useEffect } from "react";

interface ListProps {
  itemlist: {
    item: ItemProps;
    hl: {
      year: HighlightRanges | null;
      author: HighlightRanges | null;
      title: HighlightRanges | null;
    };
  }[];
  active: number;
}

const List: React.FC<ListProps> = ({ itemlist, active }) => {
  const itemRefs = useRef<(HTMLDivElement | null)[]>([]);

  useEffect(() => {
    const currentRef = itemRefs.current[active];
    if (currentRef) {
      currentRef.scrollIntoView({
        behavior: "smooth",
        block: "center",
        inline: "nearest",
      });
    }
  }, [active]);

  return (
    <div className="items">
      {itemlist.map(({ item, hl }, idx) => (
        <div
          key={idx}
          ref={(el) => {
            itemRefs.current[idx] = el;
          }}
          className={`item ${idx === active ? "active" : ""}`}
        >
          <Item item={item} hlranges={hl} />
        </div>
      ))}
    </div>
  );
};

export default List;
