import {
  Alert,
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  AlertIcon,
  Box,
  Button,
  Flex,
  Icon,
  Link,
  Text,
  Tooltip,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { BiCog, BiMessage, BiPackage, BiQuestionMark } from "react-icons/bi";
import { FiCloudLightning, FiFile, FiFolder, FiPackage } from "react-icons/fi";
import { CreatePack, OpenDirectory, SaveFile } from "../../wailsjs/go/main/App";
import { basename, dirname } from "../utils";

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
        variant="ghost"
        colorScheme="linkedin"
        leftIcon={<FiPackage />}
        ml={2}
        onClick={onOpen}
      >
        Create pack
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
                <Alert variant="left-accent" colorScheme="gray" my={2}>
                  <AlertIcon />
                  <Text>
                    Not sure what is a structured directory ? Check the{" "}
                    <Link
                      textDecor="underline"
                      href="https://github.com/olup/lunii-admin#creating-your-own-story-packs"
                    >
                      project readme
                    </Link>
                  </Text>
                </Alert>
                <Flex mt={2}>
                  <Button
                    leftIcon={<FiFolder />}
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
                      leftIcon={<FiFile />}
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
                  rightIcon={<FiPackage />}
                  colorScheme="green"
                  ref={cancelRef}
                  ml={2}
                  onClick={() =>
                    CreatePack(directoryPath, destinationPath)
                      .then(() => {
                        toast({
                          title: "Your story pack was created",
                          description: "You can now install it on your device",
                          status: "success",
                          isClosable: true,
                        });
                        handleClose();
                      })
                      .catch((err) => {
                        toast({
                          title: "Something went wrong creating the pack",
                          description: err,
                          status: "error",
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
