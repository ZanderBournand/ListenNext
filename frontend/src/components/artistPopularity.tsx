import { Badge } from "flowbite-react";
import { Flame, Star, Target, Zap } from "lucide-react";

export default function ArtistPopularity({popularity, collapse}: any) {
  let popularityWord = "";

  if (popularity >= 80) {
    popularityWord = "HOT";
  } else if (popularity >= 60) {
    popularityWord = "POPULAR";
  } else if (popularity >= 40) {
    popularityWord = "NICHE";
  } else {
    popularityWord = "NEW";
  }

  return (
    <>
      {popularity >= 80 && (
        <Badge className="ml-3 rounded-xl mt-0.5" color="purple" size="xs">
          <div className="flex flex-row items-center">
            <Flame className="h-5 w-5" />
            {!collapse && <span className="pl-1 pt-0.5">{popularityWord}</span>}
          </div>
        </Badge>
      )}
      {popularity >= 60 && popularity < 80 && (
        <Badge className="ml-3 rounded-xl mt-0.5" color="info" size="xs">
          <div className="flex flex-row items-center">
            <Star className="h-5 w-5" />
            {!collapse && <span className="pl-1 pt-0.5">{popularityWord}</span>}
          </div>
        </Badge>
      )}
      {popularity >= 40 && popularity < 60 && (
        <Badge className="ml-3 rounded-xl mt-0.5" color="light" size="xs">
          <div className="flex flex-row items-center">
            <Target className="h-5 w-5" />
            {!collapse && <span className="pl-1 pt-0.5">{popularityWord}</span>}
          </div>
        </Badge>
      )}
      {popularity < 40 && (
        <Badge className="ml-3 rounded-xl mt-0.5" color="gray" size="xs">
          <div className="flex flex-row items-center">
            <Zap className="h-5 w-5" />
            {!collapse && <span className="pl-1 pt-0.5">{popularityWord}</span>}
          </div>
        </Badge>
      )}
    </>
  );
}
