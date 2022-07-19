import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Box,
  Button,
  Flex,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import React, { useState } from "react";
import {
  MdCreate,
  MdFlashOn,
  MdOutlineFolder,
  MdOutlineInsertDriveFile,
} from "react-icons/md";
import { CreatePack, OpenDirectory, SaveFile } from "../../wailsjs/go/main/App";
import { basename, dirname } from "path-browserify";

export const NewPackModal = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef() as any;
  const toast = useToast();

  const [directoryPath, setDirectoryPath] = useState("");
  const [destinationPath, setDestinationPath] = useState("");

  const handleClose = () => {
    setDirectoryPath("");
    setDestinationPath("");
    onClose();
  };

  return (
    <>
      <Button
        variant="outline"
        colorScheme="linkedin"
        rightIcon={<MdCreate />}
        ml={2}
        onClick={onOpen}
      >
        Create story
      </Button>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={handleClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Creating a new pack
            </AlertDialogHeader>

            <AlertDialogBody>
              <Box>
                To create a new pack, first select a structured directory from
                your system
                <Flex mt={2}>
                  <Button
                    leftIcon={<MdOutlineFolder />}
                    colorScheme={directoryPath ? "green" : undefined}
                    onClick={() =>
                      OpenDirectory("Choose Directory").then(setDirectoryPath)
                    }
                  >
                    {directoryPath || "Choose directory"}
                  </Button>
                </Flex>
              </Box>
              {directoryPath && (
                <Box mt={2}>
                  Where should the pack be saved (.zip) ?
                  <Flex mt={2}>
                    <Button
                      leftIcon={<MdOutlineInsertDriveFile />}
                      colorScheme={destinationPath ? "green" : undefined}
                      onClick={() =>
                        SaveFile(
                          "Select destination",
                          dirname(directoryPath),
                          `${basename(directoryPath)}.zip`
                        ).then(setDestinationPath)
                      }
                    >
                      {destinationPath || "Select destination"}
                    </Button>
                  </Flex>
                </Box>
              )}
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={handleClose}>
                Cancel
              </Button>
              {directoryPath && destinationPath && (
                <Button
                  rightIcon={<MdFlashOn />}
                  colorScheme="green"
                  ref={cancelRef}
                  ml={2}
                  onClick={() =>
                    CreatePack(directoryPath, destinationPath).then(() => {
                      toast({
                        title: "Your story pack was created",
                        description: "You can now install it on your device",
                        status: "success",
                        duration: 9000,
                        isClosable: true,
                      });
                      handleClose();
                    })
                  }
                >
                  Create
                </Button>
              )}
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
};
