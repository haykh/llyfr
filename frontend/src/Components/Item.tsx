import { Highlight } from "@nozbe/microfuzz/react";
import { HighlightRanges } from "@nozbe/microfuzz";

import { LuScroll } from "react-icons/lu";
import { HiOutlineAcademicCap } from "react-icons/hi";
import { BsBook } from "react-icons/bs";

export interface ItemProps {
  title: string;
  author: string;
  journal: string;
  year: string;
  type: string;
  file: string;
}

const Item = ({
  item,
  hlranges,
}: {
  item: ItemProps;
  hlranges: {
    year: HighlightRanges | null;
    author: HighlightRanges | null;
    title: HighlightRanges | null;
  };
}) => (
  <p>
    {item.type.toLowerCase() === "book" ? (
      <BsBook />
    ) : item.type.toLowerCase() === "article" ? (
      <LuScroll />
    ) : item.type.toLowerCase() === "phdthesis" ? (
      <HiOutlineAcademicCap />
    ) : (
      <BsBook />
    )}
    <span className={`type ` + item.type.toLowerCase()}></span>{" "}
    <span className="year">
      <Highlight text={item.year} ranges={hlranges.year} />
    </span>{" "}
    <span className="author">
      <Highlight text={item.author} ranges={hlranges.author} />
    </span>{" "}
    <span className="journal">({item.journal})</span>{" "}
    <span className="title">
      <Highlight text={item.title} ranges={hlranges.title} />
    </span>{" "}
  </p>
);

export default Item;
