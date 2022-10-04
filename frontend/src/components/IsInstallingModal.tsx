import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Progress,
  Text,
} from "@chakra-ui/react";
import { FC, useEffect, useState } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

const totalEvents = 8;
const eventLabel: Record<string, string> = {
  STARTING: "Starting installation",
  GENERATING_BINS: "Generating binaries",
  PREPARING_ASSETS: "Preparing assets",
  CONVERTING_IMAGES: "Converting images",
  CONVERTING_AUDIOS: "Converting audios",
  WRITING_METADATA: "Writing metadata",
  COPYING: "Copying on device",
  UPDATING_INDEX: "Updating index",
  DONE: "Installation finished",
};

const getEventLabel = (eventName: string) => {
  return eventLabel[eventName] || "Processing...";
};

export const IsInstallingModal: FC<{ isOpen: boolean }> = ({ isOpen }) => {
  const [eventCount, setEventCount] = useState(0);
  const [eventName, setEventName] = useState("");
  const progress = (100 / totalEvents) * eventCount;

  useEffect(() => {
    if (isOpen) {
      // registering an event coming from the go app when performing the install
      EventsOn("INSTALL_EVENT", (event: string) => {
        setEventCount((s) => s + 1);
        setEventName(event);
      });
    } else {
      // unregistering and cleaning
      EventsOff("INSTALL_EVENT");
      setEventCount(0);
      setEventName("");
    }
  }, [isOpen]);

  return (
    <>
      <Modal isOpen={isOpen} onClose={() => {}}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Installing Pack</ModalHeader>
          <ModalBody>
            <Text mb={2}>Wait a few moments, it's almost ready.</Text>
            <Text mb={2} fontSize="sm" color="linkedin.500">
              {getEventLabel(eventName)}
            </Text>
            <Progress
              size="sm"
              value={progress}
              hasStripe
              isAnimated
              // smooth transitions
              sx={{
                "& > div:first-child": {
                  transitionProperty: "width",
                },
              }}
            />
          </ModalBody>

          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
