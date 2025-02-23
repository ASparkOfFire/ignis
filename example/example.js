// The functions in this file will ideally be coming from the SDK but for now, I left it as it is.
// bundle using ESBuild, platform=neutral

import pkg from "protobufjs";

const protobuf = pkg;

// Ensure TextEncoder is available for environments that don't support it
if (typeof TextEncoder === "undefined") {
    globalThis.TextEncoder = class {
        encode(str) {
            return new Uint8Array(str.split("").map((c) => c.charCodeAt(0)));
        }
    };
}

// Define the Protobuf schema
const root = protobuf.Root.fromJSON({
    nested: {
        FDResponse: {
            fields: {
                Body: { type: "bytes", id: 1 },
                StatusCode: { type: "int32", id: 2 },
                Length: { type: "int32", id: 3 },
                Header: { keyType: "string", type: "HeaderFields", id: 4 },
            },
        },
        HeaderFields: {
            fields: {
                fields: { rule: "repeated", type: "string", id: 1 },
            },
        },
    },
});

// Get the message types
const FDResponse = root.lookupType("FDResponse");

// Function to create a Uint8Array from a string
function stringToUint8Array(str) {
    return new Uint8Array(new TextEncoder().encode(str));
}

// Function to marshal FDResponse to binary
function marshalFDResponse(data) {
    const errMsg = FDResponse.verify(data);
    if (errMsg) {
        throw new Error(`Invalid FDResponse data: ${errMsg}`);
    }

    const message = FDResponse.create(data);
    return FDResponse.encode(message).finish();
}

// Example usage
function runExample() {
    const msg = { "msg": "Hello from Ignis JS Runtime." };
    const buf = stringToUint8Array(JSON.stringify(msg));
    const testData = {
        Body: buf,
        StatusCode: 200,
        Length: buf.length,
        Header: {
            "content-type": { fields: ["application/json"] },
            "x-custom-header": { fields: ["value1", "value2"] },
        },
    };

    const binaryData = marshalFDResponse(testData);
    const arrayBuffer = binaryData.buffer.slice(
        binaryData.byteOffset,
        binaryData.byteOffset + binaryData.byteLength
    );
    const view = new DataView(arrayBuffer);

    writebytes(view);
}

runExample();
