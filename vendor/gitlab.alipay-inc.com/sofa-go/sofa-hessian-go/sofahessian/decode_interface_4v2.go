package sofahessian

import (
	"bufio"
)

func DecodeHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	codes, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if o.maxdepth > 0 {
		o.addDepth()
		if o.depth > o.maxdepth {
			return nil, ErrDecodeMaxDepthExceeded
		}
		defer o.subDepth()
	}

	// until the go compiler will optimize it to jump table
	switch codes[0] {
	case 0:
		return DecodeStringHessian4V2(o, reader)
	case 1:
		return DecodeStringHessian4V2(o, reader)
	case 2:
		return DecodeStringHessian4V2(o, reader)
	case 3:
		return DecodeStringHessian4V2(o, reader)
	case 4:
		return DecodeStringHessian4V2(o, reader)
	case 5:
		return DecodeStringHessian4V2(o, reader)
	case 6:
		return DecodeStringHessian4V2(o, reader)
	case 7:
		return DecodeStringHessian4V2(o, reader)
	case 8:
		return DecodeStringHessian4V2(o, reader)
	case 9:
		return DecodeStringHessian4V2(o, reader)
	case 10:
		return DecodeStringHessian4V2(o, reader)
	case 11:
		return DecodeStringHessian4V2(o, reader)
	case 12:
		return DecodeStringHessian4V2(o, reader)
	case 13:
		return DecodeStringHessian4V2(o, reader)
	case 14:
		return DecodeStringHessian4V2(o, reader)
	case 15:
		return DecodeStringHessian4V2(o, reader)
	case 16:
		return DecodeStringHessian4V2(o, reader)
	case 17:
		return DecodeStringHessian4V2(o, reader)
	case 18:
		return DecodeStringHessian4V2(o, reader)
	case 19:
		return DecodeStringHessian4V2(o, reader)
	case 20:
		return DecodeStringHessian4V2(o, reader)
	case 21:
		return DecodeStringHessian4V2(o, reader)
	case 22:
		return DecodeStringHessian4V2(o, reader)
	case 23:
		return DecodeStringHessian4V2(o, reader)
	case 24:
		return DecodeStringHessian4V2(o, reader)
	case 25:
		return DecodeStringHessian4V2(o, reader)
	case 26:
		return DecodeStringHessian4V2(o, reader)
	case 27:
		return DecodeStringHessian4V2(o, reader)
	case 28:
		return DecodeStringHessian4V2(o, reader)
	case 29:
		return DecodeStringHessian4V2(o, reader)
	case 30:
		return DecodeStringHessian4V2(o, reader)
	case 31:
		return DecodeStringHessian4V2(o, reader)
	case 32:
		return DecodeBinaryHessian4V2(o, reader)
	case 33:
		return DecodeBinaryHessian4V2(o, reader)
	case 34:
		return DecodeBinaryHessian4V2(o, reader)
	case 35:
		return DecodeBinaryHessian4V2(o, reader)
	case 36:
		return DecodeBinaryHessian4V2(o, reader)
	case 37:
		return DecodeBinaryHessian4V2(o, reader)
	case 38:
		return DecodeBinaryHessian4V2(o, reader)
	case 39:
		return DecodeBinaryHessian4V2(o, reader)
	case 40:
		return DecodeBinaryHessian4V2(o, reader)
	case 41:
		return DecodeBinaryHessian4V2(o, reader)
	case 42:
		return DecodeBinaryHessian4V2(o, reader)
	case 43:
		return DecodeBinaryHessian4V2(o, reader)
	case 44:
		return DecodeBinaryHessian4V2(o, reader)
	case 45:
		return DecodeBinaryHessian4V2(o, reader)
	case 46:
		return DecodeBinaryHessian4V2(o, reader)
	case 47:
		return DecodeBinaryHessian4V2(o, reader)
	case 48:
		return DecodeStringHessian4V2(o, reader)
	case 49:
		return DecodeStringHessian4V2(o, reader)
	case 50:
		return DecodeStringHessian4V2(o, reader)
	case 51:
		return DecodeStringHessian4V2(o, reader)
	case 52:
		return DecodeBinaryHessian4V2(o, reader)
	case 53:
		return DecodeBinaryHessian4V2(o, reader)
	case 54:
		return DecodeBinaryHessian4V2(o, reader)
	case 55:
		return DecodeBinaryHessian4V2(o, reader)
	case 56:
		return DecodeInt64Hessian4V2(o, reader)
	case 57:
		return DecodeInt64Hessian4V2(o, reader)
	case 58:
		return DecodeInt64Hessian4V2(o, reader)
	case 59:
		return DecodeInt64Hessian4V2(o, reader)
	case 60:
		return DecodeInt64Hessian4V2(o, reader)
	case 61:
		return DecodeInt64Hessian4V2(o, reader)
	case 62:
		return DecodeInt64Hessian4V2(o, reader)
	case 63:
		return DecodeInt64Hessian4V2(o, reader)
	case 65:
		return DecodeBinaryHessian4V2(o, reader)
	case 66:
		return DecodeBinaryHessian4V2(o, reader)
	case 67:
		return DecodeObjectHessian4V2(o, reader)
	case 68:
		return DecodeFloat64Hessian4V2(o, reader)
	case 70:
		return DecodeBoolHessian4V2(o, reader)
	case 72:
		return DecodeMapHessian4V2(o, reader)
	case 73:
		return DecodeInt32Hessian4V2(o, reader)
	case 74:
		return DecodeDateHessian4V2(o, reader)
	case 75:
		return DecodeDateHessian4V2(o, reader)
	case 76:
		return DecodeInt64Hessian4V2(o, reader)
	case 77:
		return DecodeMapHessian4V2(o, reader)
	case 78:
		return nil, DecodeNilHessian4V2(o, reader)
	case 79:
		return DecodeObjectHessian4V2(o, reader)
	case 81:
		return DecodeRefHessian4V2(o, reader)
	case 82:
		return DecodeStringHessian4V2(o, reader)
	case 83:
		return DecodeStringHessian4V2(o, reader)
	case 84:
		return DecodeBoolHessian4V2(o, reader)
	case 86:
		return DecodeListHessian4V2(o, reader)
	case 88:
		return DecodeListHessian4V2(o, reader)
	case 89:
		return DecodeInt64Hessian4V2(o, reader)
	case 91:
		return DecodeFloat64Hessian4V2(o, reader)
	case 92:
		return DecodeFloat64Hessian4V2(o, reader)
	case 93:
		return DecodeFloat64Hessian4V2(o, reader)
	case 94:
		return DecodeFloat64Hessian4V2(o, reader)
	case 95:
		return DecodeFloat64Hessian4V2(o, reader)
	case 96:
		return DecodeObjectHessian4V2(o, reader)
	case 97:
		return DecodeObjectHessian4V2(o, reader)
	case 98:
		return DecodeObjectHessian4V2(o, reader)
	case 99:
		return DecodeObjectHessian4V2(o, reader)
	case 100:
		return DecodeObjectHessian4V2(o, reader)
	case 101:
		return DecodeObjectHessian4V2(o, reader)
	case 102:
		return DecodeObjectHessian4V2(o, reader)
	case 103:
		return DecodeObjectHessian4V2(o, reader)
	case 104:
		return DecodeObjectHessian4V2(o, reader)
	case 105:
		return DecodeObjectHessian4V2(o, reader)
	case 106:
		return DecodeObjectHessian4V2(o, reader)
	case 107:
		return DecodeObjectHessian4V2(o, reader)
	case 108:
		return DecodeObjectHessian4V2(o, reader)
	case 109:
		return DecodeObjectHessian4V2(o, reader)
	case 110:
		return DecodeObjectHessian4V2(o, reader)
	case 111:
		return DecodeObjectHessian4V2(o, reader)
	case 112:
		return DecodeListHessian4V2(o, reader)
	case 113:
		return DecodeListHessian4V2(o, reader)
	case 114:
		return DecodeListHessian4V2(o, reader)
	case 115:
		return DecodeListHessian4V2(o, reader)
	case 116:
		return DecodeListHessian4V2(o, reader)
	case 117:
		return DecodeListHessian4V2(o, reader)
	case 118:
		return DecodeListHessian4V2(o, reader)
	case 119:
		return DecodeListHessian4V2(o, reader)
	case 120:
		return DecodeListHessian4V2(o, reader)
	case 121:
		return DecodeListHessian4V2(o, reader)
	case 122:
		return DecodeListHessian4V2(o, reader)
	case 123:
		return DecodeListHessian4V2(o, reader)
	case 124:
		return DecodeListHessian4V2(o, reader)
	case 125:
		return DecodeListHessian4V2(o, reader)
	case 126:
		return DecodeListHessian4V2(o, reader)
	case 127:
		return DecodeListHessian4V2(o, reader)
	case 128:
		return DecodeInt32Hessian4V2(o, reader)
	case 129:
		return DecodeInt32Hessian4V2(o, reader)
	case 130:
		return DecodeInt32Hessian4V2(o, reader)
	case 131:
		return DecodeInt32Hessian4V2(o, reader)
	case 132:
		return DecodeInt32Hessian4V2(o, reader)
	case 133:
		return DecodeInt32Hessian4V2(o, reader)
	case 134:
		return DecodeInt32Hessian4V2(o, reader)
	case 135:
		return DecodeInt32Hessian4V2(o, reader)
	case 136:
		return DecodeInt32Hessian4V2(o, reader)
	case 137:
		return DecodeInt32Hessian4V2(o, reader)
	case 138:
		return DecodeInt32Hessian4V2(o, reader)
	case 139:
		return DecodeInt32Hessian4V2(o, reader)
	case 140:
		return DecodeInt32Hessian4V2(o, reader)
	case 141:
		return DecodeInt32Hessian4V2(o, reader)
	case 142:
		return DecodeInt32Hessian4V2(o, reader)
	case 143:
		return DecodeInt32Hessian4V2(o, reader)
	case 144:
		return DecodeInt32Hessian4V2(o, reader)
	case 145:
		return DecodeInt32Hessian4V2(o, reader)
	case 146:
		return DecodeInt32Hessian4V2(o, reader)
	case 147:
		return DecodeInt32Hessian4V2(o, reader)
	case 148:
		return DecodeInt32Hessian4V2(o, reader)
	case 149:
		return DecodeInt32Hessian4V2(o, reader)
	case 150:
		return DecodeInt32Hessian4V2(o, reader)
	case 151:
		return DecodeInt32Hessian4V2(o, reader)
	case 152:
		return DecodeInt32Hessian4V2(o, reader)
	case 153:
		return DecodeInt32Hessian4V2(o, reader)
	case 154:
		return DecodeInt32Hessian4V2(o, reader)
	case 155:
		return DecodeInt32Hessian4V2(o, reader)
	case 156:
		return DecodeInt32Hessian4V2(o, reader)
	case 157:
		return DecodeInt32Hessian4V2(o, reader)
	case 158:
		return DecodeInt32Hessian4V2(o, reader)
	case 159:
		return DecodeInt32Hessian4V2(o, reader)
	case 160:
		return DecodeInt32Hessian4V2(o, reader)
	case 161:
		return DecodeInt32Hessian4V2(o, reader)
	case 162:
		return DecodeInt32Hessian4V2(o, reader)
	case 163:
		return DecodeInt32Hessian4V2(o, reader)
	case 164:
		return DecodeInt32Hessian4V2(o, reader)
	case 165:
		return DecodeInt32Hessian4V2(o, reader)
	case 166:
		return DecodeInt32Hessian4V2(o, reader)
	case 167:
		return DecodeInt32Hessian4V2(o, reader)
	case 168:
		return DecodeInt32Hessian4V2(o, reader)
	case 169:
		return DecodeInt32Hessian4V2(o, reader)
	case 170:
		return DecodeInt32Hessian4V2(o, reader)
	case 171:
		return DecodeInt32Hessian4V2(o, reader)
	case 172:
		return DecodeInt32Hessian4V2(o, reader)
	case 173:
		return DecodeInt32Hessian4V2(o, reader)
	case 174:
		return DecodeInt32Hessian4V2(o, reader)
	case 175:
		return DecodeInt32Hessian4V2(o, reader)
	case 176:
		return DecodeInt32Hessian4V2(o, reader)
	case 177:
		return DecodeInt32Hessian4V2(o, reader)
	case 178:
		return DecodeInt32Hessian4V2(o, reader)
	case 179:
		return DecodeInt32Hessian4V2(o, reader)
	case 180:
		return DecodeInt32Hessian4V2(o, reader)
	case 181:
		return DecodeInt32Hessian4V2(o, reader)
	case 182:
		return DecodeInt32Hessian4V2(o, reader)
	case 183:
		return DecodeInt32Hessian4V2(o, reader)
	case 184:
		return DecodeInt32Hessian4V2(o, reader)
	case 185:
		return DecodeInt32Hessian4V2(o, reader)
	case 186:
		return DecodeInt32Hessian4V2(o, reader)
	case 187:
		return DecodeInt32Hessian4V2(o, reader)
	case 188:
		return DecodeInt32Hessian4V2(o, reader)
	case 189:
		return DecodeInt32Hessian4V2(o, reader)
	case 190:
		return DecodeInt32Hessian4V2(o, reader)
	case 191:
		return DecodeInt32Hessian4V2(o, reader)
	case 192:
		return DecodeInt32Hessian4V2(o, reader)
	case 193:
		return DecodeInt32Hessian4V2(o, reader)
	case 194:
		return DecodeInt32Hessian4V2(o, reader)
	case 195:
		return DecodeInt32Hessian4V2(o, reader)
	case 196:
		return DecodeInt32Hessian4V2(o, reader)
	case 197:
		return DecodeInt32Hessian4V2(o, reader)
	case 198:
		return DecodeInt32Hessian4V2(o, reader)
	case 199:
		return DecodeInt32Hessian4V2(o, reader)
	case 200:
		return DecodeInt32Hessian4V2(o, reader)
	case 201:
		return DecodeInt32Hessian4V2(o, reader)
	case 202:
		return DecodeInt32Hessian4V2(o, reader)
	case 203:
		return DecodeInt32Hessian4V2(o, reader)
	case 204:
		return DecodeInt32Hessian4V2(o, reader)
	case 205:
		return DecodeInt32Hessian4V2(o, reader)
	case 206:
		return DecodeInt32Hessian4V2(o, reader)
	case 207:
		return DecodeInt32Hessian4V2(o, reader)
	case 208:
		return DecodeInt32Hessian4V2(o, reader)
	case 209:
		return DecodeInt32Hessian4V2(o, reader)
	case 210:
		return DecodeInt32Hessian4V2(o, reader)
	case 211:
		return DecodeInt32Hessian4V2(o, reader)
	case 212:
		return DecodeInt32Hessian4V2(o, reader)
	case 213:
		return DecodeInt32Hessian4V2(o, reader)
	case 214:
		return DecodeInt32Hessian4V2(o, reader)
	case 215:
		return DecodeInt32Hessian4V2(o, reader)
	case 216:
		return DecodeInt64Hessian4V2(o, reader)
	case 217:
		return DecodeInt64Hessian4V2(o, reader)
	case 218:
		return DecodeInt64Hessian4V2(o, reader)
	case 219:
		return DecodeInt64Hessian4V2(o, reader)
	case 220:
		return DecodeInt64Hessian4V2(o, reader)
	case 221:
		return DecodeInt64Hessian4V2(o, reader)
	case 222:
		return DecodeInt64Hessian4V2(o, reader)
	case 223:
		return DecodeInt64Hessian4V2(o, reader)
	case 224:
		return DecodeInt64Hessian4V2(o, reader)
	case 225:
		return DecodeInt64Hessian4V2(o, reader)
	case 226:
		return DecodeInt64Hessian4V2(o, reader)
	case 227:
		return DecodeInt64Hessian4V2(o, reader)
	case 228:
		return DecodeInt64Hessian4V2(o, reader)
	case 229:
		return DecodeInt64Hessian4V2(o, reader)
	case 230:
		return DecodeInt64Hessian4V2(o, reader)
	case 231:
		return DecodeInt64Hessian4V2(o, reader)
	case 232:
		return DecodeInt64Hessian4V2(o, reader)
	case 233:
		return DecodeInt64Hessian4V2(o, reader)
	case 234:
		return DecodeInt64Hessian4V2(o, reader)
	case 235:
		return DecodeInt64Hessian4V2(o, reader)
	case 236:
		return DecodeInt64Hessian4V2(o, reader)
	case 237:
		return DecodeInt64Hessian4V2(o, reader)
	case 238:
		return DecodeInt64Hessian4V2(o, reader)
	case 239:
		return DecodeInt64Hessian4V2(o, reader)
	case 240:
		return DecodeInt64Hessian4V2(o, reader)
	case 241:
		return DecodeInt64Hessian4V2(o, reader)
	case 242:
		return DecodeInt64Hessian4V2(o, reader)
	case 243:
		return DecodeInt64Hessian4V2(o, reader)
	case 244:
		return DecodeInt64Hessian4V2(o, reader)
	case 245:
		return DecodeInt64Hessian4V2(o, reader)
	case 246:
		return DecodeInt64Hessian4V2(o, reader)
	case 247:
		return DecodeInt64Hessian4V2(o, reader)
	case 248:
		return DecodeInt64Hessian4V2(o, reader)
	case 249:
		return DecodeInt64Hessian4V2(o, reader)
	case 250:
		return DecodeInt64Hessian4V2(o, reader)
	case 251:
		return DecodeInt64Hessian4V2(o, reader)
	case 252:
		return DecodeInt64Hessian4V2(o, reader)
	case 253:
		return DecodeInt64Hessian4V2(o, reader)
	case 254:
		return DecodeInt64Hessian4V2(o, reader)
	case 255:
		return DecodeInt64Hessian4V2(o, reader)
	}

	return nil, ErrDecodeUnknownEncoding
}
