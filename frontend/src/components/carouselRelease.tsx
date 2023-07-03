import Carousel from "react-multi-carousel";
import "react-multi-carousel/lib/styles.css";
import RecommendationPreview from "./previewRecommendation";
import { Info } from "lucide-react";

export default function CarouselRelease({releases}: any) {
    
    const responsive = {
        desktop: {
          breakpoint: { max: 3000, min: 1024 },
          items: 3,
          slidesToSlide: 3 // optional, default to 1.
        },
        tablet: {
          breakpoint: { max: 1024, min: 464 },
          items: 2,
          slidesToSlide: 2 // optional, default to 1.
        },
        mobile: {
          breakpoint: { max: 464, min: 0 },
          items: 1,
          slidesToSlide: 1 // optional, default to 1.
        }
    };

    const ReleaseCard = ({ release }: { release: any }) => {
        return (
          <div>
            <RecommendationPreview key={release._id} release={release}/>
          </div>
        );
    };

    return (
        <div className="my-6 pr-8">
            {releases === null || releases.length == 0 && 
              <div className="flex flex-col items-center w-full py-4">
                <div className="flex flex-row items-center bg-gray-100 rounded-xl px-2 py-1">
                  <Info className="h-6 w-6 mr-2" color="#355c7d"/>
                  <span className="text-lg font-medium text-gray-600">Insufficient Data At This Time!</span>
                </div>
              </div>
            }
            <Carousel
              rewind={false}
              rewindWithAnimation={false}
              rtl={false}
              shouldResetAutoplay
              showDots={false}
              sliderClass=""
              slidesToSlide={1}
              swipeable
              additionalTransfrom={0}
              arrows
              autoPlaySpeed={3000}
              centerMode={false}
              className=""
              containerClass="container-with-dots"
              dotListClass=""
              draggable
              focusOnSelect={true}
              infinite
              itemClass=""
              keyBoardControl
              minimumTouchDrag={80}
              pauseOnHover
              renderArrowsWhenDisabled={false}
              renderButtonGroupOutside={false}
              renderDotsOutside={false}
              responsive={responsive}
            >
              {releases.map((release: any) => (
                <ReleaseCard  key={release.id} release={release}/>
              ))}
            </Carousel>
        </div>
    )
}